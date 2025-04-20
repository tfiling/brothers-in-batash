import {Shift} from '../types/shift';
import {logger} from '../utils/logger';
import {mockShifts} from '../mock/mockShifts';

// Flag to force using mock data during development
const USE_MOCK_DATA = true;

// TODO - when properly implemented, do not return mock shifts on error(current behaviour)

export interface ShiftQueryParams {
    userId?: string; // For fetching shifts for a specific user (e.g., admin view)
    startDate?: Date;
    endDate?: Date;
}

/**
 * Fetches shifts from the API.
 * If no params.userId is provided, the backend is expected to infer the user from the auth token.
 * If params.userId is provided, fetches shifts for that specific user.
 */
export async function fetchShifts(params?: ShiftQueryParams): Promise<Shift[]> {
    // Return mock data if flag is set
    if (USE_MOCK_DATA) {
        logger.info('Using mock shift data');
        let filteredShifts = [...mockShifts];
        
        // Apply client-side filtering ONLY if a specific userId is passed
        // If no userId is passed (personal view), we return all mock data, 
        // assuming backend inference will handle filtering in the real implementation.
        if (params?.userId) {
            logger.info(`Filtering mock shifts for specific user: ${params.userId}`);
            filteredShifts = mockShifts.filter(shift => {
                const isCommander = shift.commander.id === params.userId || 
                                   shift.commander.personalNumber === params.userId;
                const isAdditionalSoldier = shift.additionalSoldiers.some(
                    soldier => soldier.id === params.userId || 
                              soldier.personalNumber === params.userId
                );
                return isCommander || isAdditionalSoldier;
            });
        } else {
            logger.info('Returning all mock shifts for personal view (backend inference assumed)');
        }
        
        // Date filtering can still be applied to mock data if needed
        if (params?.startDate) {
            filteredShifts = filteredShifts.filter(shift => 
                new Date(shift.startTime) >= new Date(params.startDate!)
            );
        }
        
        if (params?.endDate) {
            filteredShifts = filteredShifts.filter(shift => 
                new Date(shift.endTime) <= new Date(params.endDate!)
            );
        }
        
        return filteredShifts;
    }

    // --- Real API Call Logic --- 
    try {
        const token = localStorage.getItem('token');

        if (!token) {
            // Return mock data temporarily if no token, but ideally should handle error/redirect
            logger.error('No token found, user must be logged in to fetch shifts');
            return mockShifts; // Or throw new Error('Authentication required');
        }

        // Build query parameters ONLY if params are provided
        const queryParams = new URLSearchParams();
        if (params?.userId) { // Only add userId if explicitly provided
            queryParams.append('userId', params.userId);
        }
        if (params?.startDate) {
            queryParams.append('startDate', params.startDate.toISOString());
        }
        if (params?.endDate) {
            queryParams.append('endDate', params.endDate.toISOString());
        }
        
        const queryString = queryParams.toString();
        // Request /api/shifts (no params if fetching for self) or /api/shifts?... (with params if provided)
        const url = `/api/shifts${queryString ? `?${queryString}` : ''}`;
        
        logger.info(`Fetching shifts from URL: ${url}`);

        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            logger.error(`Failed to fetch shifts: ${response.status} from ${url}`);
            // Return mock data temporarily, handle error appropriately in production
            return mockShifts; // Or throw new Error(`Failed to fetch shifts: ${response.statusText}`);
        }

        const shifts: Shift[] = await response.json();

        // Convert string dates to Date objects
        const processedShifts = shifts.map(shift => ({
            ...shift,
            startTime: new Date(shift.startTime),
            endTime: new Date(shift.endTime)
        }));

        logger.info(`Fetched ${shifts.length} shifts from ${url}`);
        return processedShifts;
    } catch (error) {
        logger.error('Error fetching shifts:', error);
        // Return mock data temporarily, handle error appropriately in production
        return mockShifts; // Or throw error;
    }
} 
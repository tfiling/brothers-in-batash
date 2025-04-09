import {Shift} from '../types/shift';
import {logger} from '../utils/logger';
import {mockShifts} from '../mock/mockShifts';

// Flag to force using mock data during development
const USE_MOCK_DATA = true;

/**
 * Fetches all shifts from the API
 */
export async function fetchShifts(): Promise<Shift[]> {
    // Return mock data if flag is set
    if (USE_MOCK_DATA) {
        logger.info('Using mock shift data');
        return mockShifts;
    }

    try {
        const token = localStorage.getItem('token');

        if (!token) {
            logger.error('No token found, user must be logged in to fetch shifts');
            return mockShifts; // Return mock data if no token
        }

        const response = await fetch('/api/shifts', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            logger.error(`Failed to fetch shifts: ${response.status}`);
            return mockShifts; // Return mock data on API error
        }

        const shifts: Shift[] = await response.json();

        // Convert string dates to Date objects
        const processedShifts = shifts.map(shift => ({
            ...shift,
            startTime: new Date(shift.startTime),
            endTime: new Date(shift.endTime)
        }));

        logger.info(`Fetched ${shifts.length} shifts`);
        return processedShifts;
    } catch (error) {
        logger.error('Error fetching shifts:', error);
        return mockShifts; // Return mock data on any error
    }
} 
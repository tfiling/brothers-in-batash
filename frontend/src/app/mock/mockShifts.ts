import {Shift, ShiftType, SoldierPosition} from '../types/shift';

const today = new Date();
const tomorrow = new Date(today);
tomorrow.setDate(today.getDate() + 1);

// Helper to create dates more easily
const createDate = (dayOffset: number, hours: number, minutes = 0) => {
    const date = new Date();
    date.setDate(date.getDate() + dayOffset);
    date.setHours(hours, minutes, 0, 0);
    return date;
};

export const mockShifts: Shift[] = [
    {
        id: '1',
        name: 'Morning Patrol',
        type: ShiftType.MotorizedPatrolShiftType,
        startTime: createDate(0, 8),
        endTime: createDate(0, 12),
        description: 'Morning patrol around the perimeter',
        commander: {
            id: '101',
            firstName: 'John',
            lastName: 'Doe',
            personalNumber: 'S12345',
            position: SoldierPosition.SquadCommanderPosition,
            roles: [{id: '1', name: 'Patrol Leader'}]
        },
        additionalSoldiers: [
            {
                id: '102',
                firstName: 'Jane',
                lastName: 'Smith',
                personalNumber: 'S23456',
                position: SoldierPosition.RegularSoldierPosition,
                roles: [{id: '2', name: 'Driver'}]
            }
        ]
    },
    {
        id: '2',
        name: 'Evening Watch',
        type: ShiftType.StaticPostShiftType,
        startTime: createDate(0, 18),
        endTime: createDate(0, 22),
        description: 'Evening watch at the main gate',
        commander: {
            id: '103',
            firstName: 'David',
            lastName: 'Johnson',
            personalNumber: 'S34567',
            position: SoldierPosition.CommanderPosition,
            roles: [{id: '3', name: 'Watch Commander'}]
        },
        additionalSoldiers: []
    },
    {
        id: '3',
        name: 'Night Patrol',
        type: ShiftType.MotorizedPatrolShiftType,
        startTime: createDate(0, 22),
        endTime: createDate(1, 2),
        description: 'Night patrol of the perimeter',
        commander: {
            id: '104',
            firstName: 'Michael',
            lastName: 'Brown',
            personalNumber: 'S45678',
            position: SoldierPosition.SquadLieutenantPosition,
            roles: [{id: '1', name: 'Patrol Leader'}]
        },
        additionalSoldiers: [
            {
                id: '105',
                firstName: 'Robert',
                lastName: 'Wilson',
                personalNumber: 'S56789',
                position: SoldierPosition.RegularSoldierPosition,
                roles: [{id: '2', name: 'Driver'}]
            }
        ]
    },
    {
        id: '4',
        name: 'Area Sweep',
        type: ShiftType.ProactiveOperationShiftType,
        startTime: createDate(1, 10),
        endTime: createDate(1, 14),
        description: 'Proactive sweep of designated area',
        commander: {
            id: '106',
            firstName: 'Sarah',
            lastName: 'Taylor',
            personalNumber: 'S67890',
            position: SoldierPosition.PlatoonCommanderPosition,
            roles: [{id: '4', name: 'Operations Leader'}]
        },
        additionalSoldiers: [
            {
                id: '107',
                firstName: 'James',
                lastName: 'Martin',
                personalNumber: 'S78901',
                position: SoldierPosition.RegularSoldierPosition,
                roles: [{id: '5', name: 'Scout'}]
            },
            {
                id: '108',
                firstName: 'Emily',
                lastName: 'Anderson',
                personalNumber: 'S89012',
                position: SoldierPosition.RegularSoldierPosition,
                roles: [{id: '6', name: 'Communications'}]
            }
        ]
    },
    {
        id: '5',
        name: 'Kitchen Duty',
        type: ShiftType.DailyDutyShiftType,
        startTime: createDate(1, 6),
        endTime: createDate(1, 12),
        description: 'Kitchen duty for breakfast and lunch',
        commander: {
            id: '109',
            firstName: 'Thomas',
            lastName: 'Clark',
            personalNumber: 'S90123',
            position: SoldierPosition.CommanderPosition,
            roles: [{id: '7', name: 'Duty Officer'}]
        },
        additionalSoldiers: [
            {
                id: '110',
                firstName: 'Daniel',
                lastName: 'White',
                personalNumber: 'S01234',
                position: SoldierPosition.RegularSoldierPosition,
                roles: [{id: '8', name: 'Cook'}]
            }
        ]
    }
]; 
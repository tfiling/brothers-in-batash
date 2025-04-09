export enum ShiftType {
    MotorizedPatrolShiftType,
    StaticPostShiftType,
    ProactiveOperationShiftType,
    // DailyDutyShiftType could be kitchen duty, HQ duty, etc.
    DailyDutyShiftType
}

export enum SoldierPosition {
    PlatoonCommanderPosition,
    VicePlatoonCommanderPosition,
    SquadCommanderPosition,
    SquadLieutenantPosition,
    CommanderPosition,
    RegularSoldierPosition
}

export interface SoldierRole {
    id: string;
    name: string;
    description?: string;
}

export interface Soldier {
    id: string;
    firstName: string;
    middleName?: string;
    lastName: string;
    personalNumber: string;
    position: SoldierPosition;
    roles: SoldierRole[];
}

export interface Shift {
    id: string;
    startTime: Date;
    endTime: Date;
    name: string;
    type: ShiftType;
    commander: Soldier;
    additionalSoldiers: Soldier[];
    description?: string;
    shiftTemplateId?: string;
}
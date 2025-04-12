'use client'

import {Calendar, dateFnsLocalizer, Event, SlotInfo, View} from 'react-big-calendar'
import {format} from 'date-fns/format'
import {parse} from 'date-fns/parse'
import {startOfWeek} from 'date-fns/startOfWeek'
import {getDay} from 'date-fns/getDay'
import {enUS} from 'date-fns/locale/en-US'
import 'react-big-calendar/lib/css/react-big-calendar.css'
import ProtectedRoute from './components/ProtectedRoute'
import {logger} from './utils/logger'
import {useCallback, useEffect, useState} from 'react'
import {fetchShifts} from './services/shiftService'
import {Shift, ShiftType} from './types/shift'

const locales = {
    'en-US': enUS
}

const localizer = dateFnsLocalizer({
    format,
    parse,
    startOfWeek,
    getDay,
    locales,
})

// Custom hook to manage calendar state
// Fixes a navigation bug(possibly https://github.com/jquense/react-big-calendar/issues/2720)
const useCustomCalendar = () => {
    const [view, setView] = useState<View>('week');
    const [date, setDate] = useState<Date>(new Date());

    const onView = useCallback((newView: View) => {
        logger.info('View changed to:', newView);
        setView(newView);
    }, []);

    const onNavigate = useCallback((newDate: Date) => {
        logger.info('Date navigated to:', newDate);
        setDate(newDate);
    }, []);

    return {
        view,
        date,
        onView,
        onNavigate,
    };
};

export default function HomePage() {
    const [events, setEvents] = useState<Event[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const {view, date, onView, onNavigate} = useCustomCalendar();

    const loadShifts = async () => {
        try {
            setIsLoading(true);
            const shifts = await fetchShifts();

            // Convert shifts to calendar events
            const calendarEvents = shifts.map((shift: Shift) => ({
                title: shift.name,
                start: new Date(shift.startTime),
                end: new Date(shift.endTime),
                resource: shift // Store the original shift data in the resource field
            }));

            setEvents(calendarEvents);
            setError(null);
        } catch (err) {
            setError('Failed to load shifts');
            logger.error('Error loading shifts:', err);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        logger.info('HomePage component mounted')
        loadShifts();
        
        return () => {
            logger.info('HomePage component unmounted')
        }
    }, [])

    const handleSelectEvent = (event: Event): void => {
        const shift = event.resource as Shift;
        logger.info('Shift selected:', shift);

        // You could display a modal with shift details here
        alert(`Shift: ${shift.name}\nCommander: ${shift.commander.firstName} ${shift.commander.lastName}\nStart: ${shift.startTime.toLocaleString()}\nEnd: ${shift.endTime.toLocaleString()}`);
    }

    const handleSelectSlot = (slotInfo: SlotInfo): void => {
        logger.info('Time slot selected:', slotInfo)
    }

    const handleRefresh = () => {
        loadShifts();
    };

    // Get class name for shift type
    const eventPropGetter = (event: Event) => {
        const shift = event.resource as Shift;
        let className = '';

        if (shift) {
            switch (shift.type) {
                case ShiftType.MotorizedPatrolShiftType:
                    className = 'shift-motorized-patrol';
                    break;
                case ShiftType.StaticPostShiftType:
                    className = 'shift-static-post';
                    break;
                case ShiftType.ProactiveOperationShiftType:
                    className = 'shift-proactive-operation';
                    break;
                case ShiftType.DailyDutyShiftType:
                    className = 'shift-daily-duty';
                    break;
            }
        }

        return {className};
    };

    return (
        <ProtectedRoute>
            <div className="calendar-container">
                <div className="calendar-header">
                    <h1>Shift Calendar</h1>
                    <button
                        className="refresh-button"
                        onClick={handleRefresh}
                        disabled={isLoading}
                    >
                        {isLoading ? 'Loading...' : 'Refresh'}
                    </button>
                </div>
                {isLoading && <div className="loading">Loading shifts...</div>}
                {error && <div className="error">{error}</div>}
                <Calendar
                    localizer={localizer}
                    events={events}
                    startAccessor="start"
                    endAccessor="end"
                    style={{height: '100%'}}
                    onSelectEvent={handleSelectEvent}
                    onSelectSlot={handleSelectSlot}
                    eventPropGetter={eventPropGetter}
                    selectable
                    view={view}
                    date={date}
                    onView={onView}
                    onNavigate={onNavigate}
                    className="calendar-light-theme"
                />
            </div>
        </ProtectedRoute>
    )
} 
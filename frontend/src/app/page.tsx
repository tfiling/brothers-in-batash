'use client'

import {Calendar, dateFnsLocalizer, Event, SlotInfo} from 'react-big-calendar'
import {format} from 'date-fns/format'
import {parse} from 'date-fns/parse'
import {startOfWeek} from 'date-fns/startOfWeek'
import {getDay} from 'date-fns/getDay'
import {enUS} from 'date-fns/locale/en-US'
import 'react-big-calendar/lib/css/react-big-calendar.css'
import ProtectedRoute from './components/ProtectedRoute'
import {logger} from './utils/logger'
import {useEffect} from 'react'

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

const events = [
    {
        title: 'Team Meeting',
        start: new Date(2024, 2, 20, 10, 0),
        end: new Date(2024, 2, 20, 11, 30),
    },
    {
        title: 'Lunch Break',
        start: new Date(2024, 2, 20, 12, 0),
        end: new Date(2024, 2, 20, 13, 0),
    },
    {
        title: 'Project Deadline',
        start: new Date(2024, 2, 22, 9, 0),
        end: new Date(2024, 2, 22, 17, 0),
    },
]

export default function HomePage() {
    useEffect(() => {
        logger.info('HomePage component mounted')
        return () => {
            logger.info('HomePage component unmounted')
        }
    }, [])

    const handleSelectEvent = (event: Event): void => {
        logger.info('Event selected:', event.title)
    }

    const handleSelectSlot = (slotInfo: SlotInfo): void => {
        logger.info('Time slot selected:', slotInfo)
    }

    return (
        <ProtectedRoute>
            <div className="calendar-container">
                <Calendar
                    localizer={localizer}
                    events={events}
                    startAccessor="start"
                    endAccessor="end"
                    style={{height: '100%'}}
                    onSelectEvent={handleSelectEvent}
                    onSelectSlot={handleSelectSlot}
                />
            </div>
        </ProtectedRoute>
    )
} 
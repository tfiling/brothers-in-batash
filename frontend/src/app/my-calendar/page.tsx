'use client'

import ProtectedRoute from '../components/ProtectedRoute'
import CalendarView from '../components/CalendarView'

export default function MyCalendarPage() {
    return (
        <ProtectedRoute>
            <CalendarView 
                mode="personal" 
                title="My Shifts" 
            />
        </ProtectedRoute>
    )
} 
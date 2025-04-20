'use client'

import ProtectedRoute from './components/ProtectedRoute'
import CalendarView from './components/CalendarView'

export default function HomePage() {
    return (
        <ProtectedRoute>
            <CalendarView mode="global" title="Company Shift Calendar" />
        </ProtectedRoute>
    )
} 
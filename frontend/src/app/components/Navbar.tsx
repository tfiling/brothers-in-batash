'use client'

import {useAuth} from '../context/AuthContext'
import Link from 'next/link'
import {usePathname} from 'next/navigation'

export default function Navbar() {
    const {user, logout} = useAuth()
    const pathname = usePathname()

    // Determine which nav link is active
    const isActive = (path: string) => {
        return pathname === path ? 'navbar-link-active' : ''
    }

    return (
        <nav className="navbar">
            <div className="navbar-container">
                <div className="navbar-brand">Squad Mate</div>
                {user && (
                    <div className="navbar-links">
                        <Link href="/" className={`navbar-link ${isActive('/')}`}>
                            Company Calendar
                        </Link>
                        <Link href="/my-calendar" className={`navbar-link ${isActive('/my-calendar')}`}>
                            My Shifts
                        </Link>
                    </div>
                )}
                <div className="navbar-user">
                    {user && (
                        <>
                            <span className="navbar-username">Welcome, {user}</span>
                            <button
                                onClick={logout}
                                className="navbar-button"
                            >
                                Logout
                            </button>
                        </>
                    )}
                </div>
            </div>
        </nav>
    )
} 
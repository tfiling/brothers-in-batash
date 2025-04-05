'use client'

import {useAuth} from '../context/AuthContext'

export default function Navbar() {
    const {user, logout} = useAuth()

    return (
        <nav className="navbar">
            <div className="navbar-container">
                <div className="navbar-brand">Squad Mate</div>
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
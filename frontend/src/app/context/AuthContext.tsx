'use client'

import {createContext, ReactNode, useContext, useEffect, useState} from 'react'
import {useRouter} from 'next/navigation'
import {logger} from '../utils/logger'

interface AuthContextType {
    user: string | null
    login: (username: string, password: string) => Promise<boolean>
    register: (username: string, password: string) => Promise<boolean>
    logout: () => void
    error: string | null
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({children}: { children: ReactNode }) {
    const [user, setUser] = useState<string | null>(null)
    const [error, setError] = useState<string | null>(null)
    const router = useRouter()

    useEffect(() => {
        // Check if user is logged in on mount
        const storedUser = localStorage.getItem('user')
        if (storedUser) {
            setUser(storedUser)
        }
    }, [])

    const login = async (username: string, password: string): Promise<boolean> => {
        try {
            logger.info('Making login request...')
            setError(null)
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({username, password}),
            })
            logger.info('request url:', response.url)
            console.log('Response status:', response.status)

            if (!response.ok) {
                const data = await response.json().catch(() => ({}))
                console.log('Error response:', data)
                setError(data.error || 'Login failed')
                return false
            }

            const data = await response.json()
            logger.info('login response body:', data)
            setUser(data.username)
            localStorage.setItem('user', data.username)
            return true
        } catch (error) {
            console.error('Login error:', error)
            setError('An error occurred during login')
            return false
        }
    }

    const register = async (username: string, password: string): Promise<boolean> => {
        try {
            setError(null)
            const response = await fetch('/api/v1/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({username, password}),
            })

            if (!response.ok) {
                const data = await response.json()
                setError(data.error || 'Registration failed')
                return false
            }

            return true
        } catch {
            setError('An error occurred during registration')
            return false
        }
    }

    const logout = () => {
        setUser(null)
        localStorage.removeItem('user')
        router.push('/login')
    }

    return (
        <AuthContext.Provider value={{user, login, register, logout, error}}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    logger.info('useAuth hook called(logger)')
    const context = useContext(AuthContext)
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}
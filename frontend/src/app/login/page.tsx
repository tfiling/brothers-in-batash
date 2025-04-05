'use client'

import {FormEvent, useState} from 'react'
import {useRouter} from 'next/navigation'
import Link from 'next/link'
import {useAuth} from '../context/AuthContext'

export default function LoginPage() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [isLoading, setIsLoading] = useState(false)
    const {login, error} = useAuth()
    const router = useRouter()

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        setIsLoading(true)

        const success = await login(username, password)
        setIsLoading(false)

        if (success) {
            router.push('/')
        }
    }

    return (
        <div className="form-container">
            <div className="form-card">
                <h2 className="form-title">Sign in to your account</h2>
                <form className="form" onSubmit={handleSubmit}>
                    <div className="form-group">
                        <input
                            id="username"
                            name="username"
                            type="text"
                            required
                            className="form-input"
                            placeholder="Username"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                        />
                    </div>
                    <div className="form-group">
                        <input
                            id="password"
                            name="password"
                            type="password"
                            required
                            className="form-input"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>

                    {error && (
                        <div className="form-error">{error}</div>
                    )}

                    <button
                        type="submit"
                        disabled={isLoading}
                        className="form-button"
                    >
                        {isLoading ? 'Signing in...' : 'Sign in'}
                    </button>

                    <Link href="/register" className="form-link">
                        Don&apos;t have an account? Register
                    </Link>
                </form>
            </div>
        </div>
    )
} 
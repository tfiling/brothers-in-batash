'use client'

import {FormEvent, useState} from 'react'
import {useRouter} from 'next/navigation'
import Link from 'next/link'
import {useAuth} from '../context/AuthContext'

export default function RegisterPage() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [confirmPassword, setConfirmPassword] = useState('')
    const [isLoading, setIsLoading] = useState(false)
    const [validationError, setValidationError] = useState<string | null>(null)
    const {register, error} = useAuth()
    const router = useRouter()

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        setValidationError(null)

        if (password !== confirmPassword) {
            setValidationError('Passwords do not match')
            return
        }

        if (password.length < 4) {
            setValidationError('Password must be at least 4 characters long')
            return
        }

        if (username.length < 4) {
            setValidationError('Username must be at least 4 characters long')
            return
        }

        setIsLoading(true)

        const success = await register(username, password)
        setIsLoading(false)

        if (success) {
            router.push('/login')
        }
    }

    return (
        <div className="form-container">
            <div className="form-card">
                <h2 className="form-title">Create your account</h2>
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
                    <div className="form-group">
                        <input
                            id="confirmPassword"
                            name="confirmPassword"
                            type="password"
                            required
                            className="form-input"
                            placeholder="Confirm Password"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                        />
                    </div>

                    {(error || validationError) && (
                        <div className="form-error">
                            {validationError || error}
                        </div>
                    )}

                    <button
                        type="submit"
                        disabled={isLoading}
                        className="form-button"
                    >
                        {isLoading ? 'Creating account...' : 'Create account'}
                    </button>

                    <Link href="/login" className="form-link">
                        Already have an account? Sign in
                    </Link>
                </form>
            </div>
        </div>
    )
} 
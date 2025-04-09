'use client'
// TODO: consider relocating this file
import {useEffect} from 'react'

export default function ClientLogger() {
    useEffect(() => {
        console.log('Client component mounted')  // This will show only in browser console
    }, [])

    return null
} 
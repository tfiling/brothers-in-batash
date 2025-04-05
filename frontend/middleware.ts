import type {NextRequest} from 'next/server'
import {NextResponse} from 'next/server'

export function middleware(request: NextRequest) {
    console.log('Middleware processing request:', request.url)  // This will show in server logs
    return NextResponse.next()
}

export const config = {
    matcher: '/api/:path*',
} 
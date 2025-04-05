import {logger} from '@/app/utils/logger'
import type {NextConfig} from 'next'

const nextConfig: NextConfig = {
    reactStrictMode: true,
    // output: 'standalone',  // Optimizes for Docker(in prod mode?)
    async rewrites() {
        const apiUrl = process.env.API_URL || 'http://localhost:3000'
        console.log('API URL:', apiUrl) // Add this for debugging
        logger.info('rewrites called')
        return [
            {
                source: '/api/:path*',
                destination: `${apiUrl}/api/v1/:path*`
            }
        ]
    }
}

export default nextConfig
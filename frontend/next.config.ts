import type {NextConfig} from 'next'

const nextConfig: NextConfig = {
    reactStrictMode: true,
    output: 'standalone',  // Optimizes for Docker
    async rewrites() {
        return [
            {
                source: '/api/:path*',
                destination: (process.env.API_URL || 'http://localhost:3000') + '/:path*' // Proxy API requests to backend
            }
        ]
    }
}

export default nextConfig
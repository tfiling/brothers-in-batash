import './globals.css'
import {AuthProvider} from './context/AuthContext'
import Navbar from './components/Navbar'

console.log('Next.js server initialized application')

export default function RootLayout({
                                       children,
                                   }: {
    children: React.ReactNode
}) {
    return (
        <html lang="en">
        <body suppressHydrationWarning={true}>
        <AuthProvider>
            <Navbar/>
            {children}
        </AuthProvider>
        </body>
        </html>
    )
} 
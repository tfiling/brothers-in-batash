type LogLevel = 'info' | 'warn' | 'error'
type LogArgs = unknown[]

const formatMessage = (level: LogLevel, message: string): string => {
    const timestamp = new Date().toISOString()
    return `[${timestamp}] [${level.toUpperCase()}] ${message}`
}

export const logger = {
    info: (message: string, ...args: LogArgs): void => {
        console.log(formatMessage('info', message), ...args)
    },
    warn: (message: string, ...args: LogArgs): void => {
        console.warn(formatMessage('warn', message), ...args)
    },
    error: (message: string, ...args: LogArgs): void => {
        console.error(formatMessage('error', message), ...args)
    }
} 
import './globals.css'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Chinese Learning Graph',
  description: 'Interactive Chinese word learning through visual connections',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: process.env.NODE_ENV === 'production' 
          ? 'http://backend:8080/:path*'
          : 'http://localhost:8080/:path*'
      }
    ]
  }
}

module.exports = nextConfig
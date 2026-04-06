/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  async rewrites() {
    return [
      {
        source: '/graphql',
        destination: process.env.NODE_ENV === 'production' 
          ? 'http://backend:8080/graphql'
          : 'http://localhost:8080/graphql'
      },
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
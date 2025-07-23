/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async rewrites() {
    // If backend URL is configured, proxy API requests to external backend
    if (process.env.NEXT_PUBLIC_API_BASE_URL) {
      return [
        {
          source: '/api/:path*',
          destination: `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/:path*`,
        },
      ]
    }
    
    // Default behavior: use internal API handlers
    return [
      {
        source: '/api/:path*',
        destination: '/api/:path*',
      },
    ]
  },
  output: "standalone"
}

export default nextConfig

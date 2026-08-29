import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        hostname: 'localhost',
      },
    ],
    deviceSizes: [320, 1320, 1920, 2560],
    formats: ['image/avif', 'image/webp'],
  },

};

export default nextConfig;
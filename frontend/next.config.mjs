// Set up a proxy for API requests to the Go backend
const nextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:8080/api/:path*", // Proxy to Go backend
      },
    ];
  },
};

export default nextConfig;

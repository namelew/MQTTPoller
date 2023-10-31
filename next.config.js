/**
 * @type {import('next').NextConfig}
*/

const nextConfig = {
  "output": "standalone",
  "rewrites": async () => {
    return [
      {
        "source": "/api/:path*",
        "destination": `${process.env.ORQUESTRATOR_ADRESS ? process.env.ORQUESTRATOR_ADRESS : 'http://localhost:8000'}/orquestrator/:path*`,
      },
    ];
  }
}

module.exports = nextConfig

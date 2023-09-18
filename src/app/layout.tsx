import { Providers } from "./_providers"
import AppDefault from "./components/default"
import Footer from "./components/footer"
import Navbar from "./components/navbar"

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>
        <Providers>
          <Navbar />
            <AppDefault>
              {children}
            </AppDefault>
          <Footer />
        </Providers>
      </body>
    </html>
  )
}

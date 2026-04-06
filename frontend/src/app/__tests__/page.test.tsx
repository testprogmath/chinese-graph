import { render, screen } from '@testing-library/react'
import Home from '../page'

// Mock Next.js router
jest.mock('next/navigation', () => ({
  useRouter() {
    return {
      push: jest.fn(),
      replace: jest.fn(),
      prefetch: jest.fn()
    }
  },
  useSearchParams() {
    return new URLSearchParams()
  },
  usePathname() {
    return ''
  }
}))

describe('Home page', () => {
  it('renders the main heading', () => {
    render(<Home />)
    
    const heading = screen.getByRole('heading', {
      name: /chinese learning graph/i,
    })
    
    expect(heading).toBeInTheDocument()
  })

  it('displays the description text', () => {
    render(<Home />)
    
    const description = screen.getByText(/learn chinese words through visual connections/i)
    
    expect(description).toBeInTheDocument()
  })

  it('shows development status', () => {
    render(<Home />)
    
    const status = screen.getByText(/frontend under development/i)
    
    expect(status).toBeInTheDocument()
  })
})
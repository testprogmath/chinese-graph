import { render, screen } from '@testing-library/react'
import { MockedProvider } from '@apollo/client/testing'
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

const renderWithApollo = (component: React.ReactElement) => {
  return render(
    <MockedProvider mocks={[]}>
      {component}
    </MockedProvider>
  )
}

describe('Home page', () => {
  it('renders the main heading', () => {
    renderWithApollo(<Home />)
    
    const heading = screen.getByRole('heading', {
      name: /chinese learning graph/i,
    })
    
    expect(heading).toBeInTheDocument()
  })

  it('displays the description text', () => {
    renderWithApollo(<Home />)
    
    const description = screen.getByText(/learn chinese words through visual connections/i)
    
    expect(description).toBeInTheDocument()
  })

  it('shows the footer text', () => {
    renderWithApollo(<Home />)
    
    const footer = screen.getByText(/explore word relationships.*practice with context.*learn through connections/i)
    
    expect(footer).toBeInTheDocument()
  })
})
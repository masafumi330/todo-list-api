import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Home from './pages/Home'
import SignUp from './pages/SignUp'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/sign_up" element={<SignUp />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App

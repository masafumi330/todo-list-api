import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

function SignUp() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const navigate = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    
    try {
      const response = await fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      })

      if (response.ok) {
        // 登録成功
        alert('登録が完了しました')
        navigate('/')
      } else {
        // エラー処理
        alert('登録に失敗しました')
      }
    } catch (error) {
      console.error('Error:', error)
      alert('エラーが発生しました')
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-orange-50 via-blue-50 to-purple-50 p-4 sm:p-8">
      <main className="w-full max-w-lg bg-white rounded-3xl shadow-2xl p-8 sm:p-12">
        <h1 className="text-3xl sm:text-4xl font-bold text-center text-slate-900 mb-8">
          新規登録
        </h1>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-slate-700 mb-2">
              メールアドレス
            </label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="w-full px-4 py-3 rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
              placeholder="example@example.com"
            />
          </div>
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-slate-700 mb-2">
              パスワード
            </label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full px-4 py-3 rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
              placeholder="パスワードを入力"
            />
          </div>
          <button
            type="submit"
            className="w-full px-6 py-3 rounded-full font-semibold text-white bg-blue-600 hover:bg-blue-700 shadow-lg shadow-blue-600/30 hover:-translate-y-0.5 transition-all duration-200 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
          >
            登録する
          </button>
        </form>
      </main>
    </div>
  )
}

export default SignUp

import { Link } from 'react-router-dom'

function Home() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-orange-50 via-blue-50 to-purple-50 p-4 sm:p-8">
      <main className="w-full max-w-lg bg-white rounded-3xl shadow-2xl p-8 sm:p-12">
        <h1 className="text-3xl sm:text-4xl font-bold text-center text-slate-900 mb-8">
          ToDoアプリです
        </h1>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link
            to="/sign_up"
            className="min-w-[140px] px-6 py-3 rounded-full font-semibold text-white bg-blue-600 hover:bg-blue-700 shadow-lg shadow-blue-600/30 hover:-translate-y-0.5 transition-all duration-200 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 text-center"
          >
            新規登録
          </Link>
          <button
            type="button"
            className="min-w-[140px] px-6 py-3 rounded-full font-semibold text-blue-600 bg-white border-2 border-blue-600 hover:bg-blue-50 shadow-lg shadow-slate-900/10 hover:-translate-y-0.5 transition-all duration-200 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
          >
            ログイン
          </button>
        </div>
      </main>
    </div>
  )
}

export default Home

import { use, useState } from 'react';
import { useInvoiceStore } from '@/store/useInvoiceStore';
import { useRouter } from 'next/router';
import api from '@/utils/api';
import Cookies from 'js-cookie';

export default function loginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const setSAuth = useInvoiceStore((state) => state.setAuth);
  const router = useRouter();

  const handleLogin = async (e: React.FormEvent) =>  {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const response = await api.post('login', {
        username,
        password,
      });

      const { token, role } = response.data;
      setSAuth( token, role );
      Cookies.set('token', token, { expires: 1});
      alert(`login berhasil sebagai ${role}`);
      
      router.push('/wizard');      
    } catch (err: any) {
      setError(err.response?.data?.message || 'Gagal terhubung ke server');
    } finally {
      setIsLoading(false);
    }
  }
  
  return (
    <div className='flex min-h-screen items-center justify-center bg-slate-900 text-white'>
    <div className='w-full max-w-md p-8 bg-slatte-800 rounded-xl shadow-2xl border border-slate-700'>
      <h1 className='text-3xl font-bold mb-w text-center text-blue-400'>Fleetfy</h1>
      <p className='text-slate-400 text-center mb-8'>Logistic& Fleet Management </p>
      {error && (
        <div className='bg-red-500/10 border border-red-500 text-red-500 p-3 rounded mb-4 text-sm'>
          {error}
        </div> 
      )}

      <form onSubmit={handleLogin} className='space-y-4'>
        <div>
          <label className='block text-sm font-medium mb-1'>Username</label>
          <input
          type="text"
          required
          className='w-full p-2.5 bg-slate-700 border border-slate-600 rounded focus:ring-2 focus:ring-blue-500 outlline-none transition'
          placeholder='admin / kerani'
          value={username}
          onChange={(e) => setUsername(e.target.value)}/>
        </div>
        <div>
          <label className='block text-sm font-medium mb-1'>Password</label>
          <input type="password" required
          className='w-full p-2.5 bg-slate-700 border border-slate-600 rounded focus:ring-2 focus:ring-blue-500 outline-none transition'
          placeholder='••••••••'
          value={password}
          onChange={(e) => setPassword(e.target.value)}/>
        </div>
        <button type='submit' disabled={isLoading} className='w-full bg-blue-600 hiver:bg-blue-700 text-white font-bold py-2.5 rounded transition duration-200 disabled:opacity-50'>
          {isLoading ? 'Authenticating...' : 'Sign In'}
        </button>
      </form>
        <div className='mt-6 text-center text-xs text-slate-500'>
          <p>Admin: admin / admin123</p>
          <p>Kerani: kerani / kerani123</p>
        </div>
      </div>
    </div>
  );
};
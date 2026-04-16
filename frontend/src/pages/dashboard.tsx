import { useEffect, useState } from 'react';
import api from '@/utils/api';
import { useRouter } from 'next/router';

export default function Dashboard() {
  const [invoices, setInvoices] = useState([]);
  const router = useRouter();

  useEffect(() => {
    api.get('/invoices')
      .then(res => setInvoices(res.data))
      .catch(err => console.error("Gagal ambil data", err));
  }, []);

  return (
    <div className="min-h-screen bg-slate-900 text-white p-6 md:p-10">
      <div className="max-w-5xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-blue-400">Invoice Dashboard</h1>
          <button 
            onClick={() => router.push('/wizard')}
            className="bg-blue-600 hover:bg-blue-700 px-5 py-2 rounded-lg font-bold transition"
          >
            + Buat Invoice Baru
          </button>
        </div>

        <div className="grid gap-6">
          {invoices.length === 0 ? (
            <p className="text-slate-500 text-center py-20">Belum ada invoice yang dibuat.</p>
          ) : (
            invoices.map((inv: any) => (
              <div key={inv.id} className="bg-slate-800 border border-slate-700 rounded-xl p-6 shadow-lg">
                <div className="flex justify-between items-start border-b border-slate-700 pb-4 mb-4">
                  <div>
                    <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Nomor Invoice</p>
                    <p className="text-xl font-mono text-blue-300">{inv.invoice_number}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-xs text-slate-500 font-bold uppercase tracking-wider">Total Tagihan</p>
                    <p className="text-xl font-bold text-green-400">Rp{inv.total_amount?.toLocaleString()}</p>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4 text-sm mb-4">
                  <div>
                    <span className="text-slate-500">Pengirim:</span>
                    <p className="font-semibold">{inv.sender_name}</p>
                  </div>
                  <div>
                    <span className="text-slate-500">Penerima:</span>
                    <p className="font-semibold">{inv.receiver_name}</p>
                  </div>
                </div>

                <div className="bg-slate-900/50 rounded-lg p-3">
                  <p className="text-xs font-bold text-slate-600 mb-2 uppercase">Detail Barang:</p>
                  {inv.details?.map((item: any, idx: number) => (
                    <div key={idx} className="flex justify-between text-xs text-slate-400 mb-1">
                      <span>{item.name} (x{item.quantity})</span>
                      <span>Rp{item.subtotal?.toLocaleString()}</span>
                    </div>
                  ))}
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
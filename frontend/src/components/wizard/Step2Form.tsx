import { useState, useEffect } from 'react';
import { useInvoiceStore } from '@/store/useInvoiceStore';
import api from '@/utils/api';

export default function Step2Form({ onNext, onPrev}: { onNext: () => void; onPrev: () => void}) {
    const { items, addItem} = useInvoiceStore();
    const [code, setCode] = useState('');
    const [qty, setQty] = useState(1);
    const [loading, setLoading] = useState(false);
    const [foundItem, setFoundItem] = useState<any>(null);
    const [error, setError] = useState('');

    useEffect(() => {
        if (!code) {
            setFoundItem(null);
            return;
        }

        const controller = new AbortController();
        const timeoutId = setTimeout(async () => {
            setLoading(true);
            try {
                const res = await api.get(`items?code=${code}`, {
                    signal: controller.signal
                });
                setFoundItem(res.data);
                setError('');
            } catch (err: any) {
                if (err.name !== "CancelError" ) {
                    setFoundItem(null);
                    setError('Barang tidak ditemukan!')
                }
            } finally {
                setLoading(false);
            }

        }, 500);

        return () => {
            clearTimeout(timeoutId);
            controller.abort();
        };
    }, [code]);

    const handleAddItem = () => {
        if (foundItem) { 
            addItem({
                ...foundItem, 
                quantity: qty,
                subtotal: foundItem.price * qty
            });
            setCode('');
            setQty(1);
            setFoundItem(null);
        }
    };

    return (
        <div className='space-y-6'>
            <h2 className='text-xl font-bold border-b border-slate-700 pb-2'>Step 2: Tambah Barang</h2>
            <div className='flex gap-4 items-end bg-slate-700/50 p-4 rounded-lg'>
                <div className='flex-1'>
                    <label className='block text-sm font-medium text-slate-300 mb-1'>Kode Barang</label>
                    <input value={code} onChange={(e) => setCode(e.target.value)} className="w-full bg-slate-800 border border-slate-600 p-2 rounded outline-none focus:border-blue-500" placeholder='Ketik kode...' />
                </div>
                <div className='w-24'>
                    <label className='block text-sm mb-1 text-slate-400'>Qty</label>
                    <input type="number" value={qty} onChange={(e) => setQty(Number(e.target.value))} className='w-full bg-slate-800 border border-slate-600 p-2 rounded outline-none' />
                </div>
                <button onClick={handleAddItem} disabled={!foundItem || loading} className='bg-green-600 px-4 py-2 rounded font-bold hover:bg-green-700 disabled:opacity-50 h-[42px]'>
                    Tambah
                </button>
            </div>
            {loading && <p className='text-blue-400 text-sm animate-pulse'>Mencari barang...</p>}
            {error && !loading && (
                <p className='text-red-400 text-sm font-medium bg-red-400/10 p-2 rounded border border-red-400/20'>
                    ⚠️ {error}
                </p>
            )}
            {foundItem && <p className='text-green-400 text-sm'>Terdeteksi: {foundItem.name} - Rp{foundItem.price.toLocaleString()}</p>}

            <table className='w-full text-left text-sm'>
                <thead>
                    <tr className='Border-b border-slate-700 text-slate-400'>
                        <th className='py-2'>Nama Barang</th>
                        <th className='py-2'>Harga</th>
                        <th className='py-2'>Qty</th>
                        <th className='py-2 text-right'>Subtotal</th>
                    </tr>
                </thead>
                <tbody>
                    {items.map((item, idx) => (
                        <tr key={idx} className='border-b-slate-700/50'>
                            <td className='py-2'>{item.name}</td>
                            <td className='py-2'>Rp{item.price.toLocaleString()}</td>
                            <td className='py-2'>{item.quantity}</td>
                            <td className='py-2 text-right'>Rp{item.subtotal.toLocaleString()}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <div className='flex justify-between pt-4'>
                <button onClick={onPrev} className='text-slate-400 hover:text-white'>Kembali</button>
                <button onClick={onNext} disabled={items.length === 0} className='bg-blue-600 px-6 py-2 rounded font-bold hover:bg-blue-700 disabled:opacity-50'>
                    Review Invoice
                </button>  
            </div>
        </div>
    );
}
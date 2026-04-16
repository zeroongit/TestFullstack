import { useInvoiceStore } from "@/store/useInvoiceStore";
import api from "@/utils/api";
import { useRouter } from "next/router";

export default function Step3Review ({onPrev }: { onPrev: () => void }) {
    const { step1, resetForm, items } = useInvoiceStore();
    const router = useRouter();

    const total = items.reduce((acc, curr) => acc + curr.subtotal, 0);
    const handleSubmit = async () => {
        try {
            await api.post('/invoices', {
                ...step1,
                items: items,
                total_amount: total
            });
            alert('Invoice berhasil disimpan!');
            resetForm();
            router.push('/');
        } catch (err) {
            alert('Gagal menyimpan invoice. Coba lagi.');
        }
    };

    return (
        <div className="space-y-6">
            <h2 className="text-xl font-bold border-b border-slate-700 pb-2">Step 3: Konfirmasi Final</h2>
            <div className="grid grid-cols-2 gap-8 text-sm">
                <div className="space-y-1">
                    <p className="text-slate-500 uppercase text-xs font-bold">Pengirim</p>
                    <p className="font-bold">{step1.senderName}</p>
                    <p className="text-slate-400">{step1.senderAddress}</p>
                </div>
                <div className="space-y-1">
                    <p className="text-slate-500 uppercase text-xs font-bold">Penerima</p>
                    <p className="font-bold">{step1.receiverName}</p>
                    <p className="text-slate-400">{step1.receiverAddress}</p>
                </div>
            </div>

            <div className="bg-slate-900/50 p-4 rounded-lg">
                <p className="text-slate-500 uppercase text-xs font-bold mb-2">Rincian Barang</p>
                {items.map((item, idx) => (
                    <div key={idx} className="flex justify-between text-sm mb-1">
                        <span>{item.name} (x{item.quantity})</span>
                        <span>Rp{item.subtotal.toLocaleString()}</span>
                    </div>
                ))}
                <div className="border-t border-slate-700 mt-2 pt-2 flex justify-between font-bold text-lg text-blue-400">
                    <span>Total Tagihan</span>
                    <span>Rp{total.toLocaleString()}</span>
                </div>
            </div>
            <div className="flex justify-between pt-4">
                <button onClick={onPrev} className="text-slate-400 hover:text-white">Cek  Ulang Barang</button>
                <button onClick={handleSubmit} className="bg-green-600 px-8 py-3 rounded-lg font-bold hover:bg-green-700 shadow-lg shadow-green-900/20">
                    Simpan & Selesai
                </button>
            </div>
        </div>
    )
};
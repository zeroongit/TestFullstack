import { useInvoiceStore } from "@/store/useInvoiceStore";

export default function Step1Form({ onNext }: { onNext: () => void}) {
    const { step1, updateStep1 } = useInvoiceStore();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        updateStep1({...step1, [e.target.name]: e.currentTarget.value});
    };

    return (
        <div className="space-y-4">
            <h2 className="text-xl font-bold border-b border-slate-700 pb-2">Step 1: Data Pengiriman</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm mb-1 text-slate-400">Nama Pengirim</label>
                    <input name="senderName" value={step1.senderName} onChange={handleChange} className="w-full bg-slate-700 border border-slate-600 p-2 rounded outline-none focus:border-blue-500" />
                </div>
                <div>
                    <label className="block text-sm mb-1 text-slate-400">Nama Penerima</label>
                    <input name="receiverName" value={step1.receiverName} onChange={handleChange} className="w-full bg-slate-700 border border-slate-600 p-2 rounded outline-none focus:border-blue-500" />
                </div>
            </div>
            <div>
                <label className="block text-sm mb-1 text-slate-400">Alamat Pengirim</label>
                <textarea name="senderAddress" value={step1.senderAddress} onChange={handleChange} className="w-full bg-slate-700 border border-slate-600 p-2 rounded h-20 outline-none focus:border-blue-500" />
            </div>
            <div>
                <label className="block text-sm mb-1 text-slate-400">Alamat Penerima</label>
                <textarea name="receiverAddress" value={step1.receiverAddress} onChange={handleChange} className="w-full bg-slate-700 border border-slate-600 p-2 rounded h-20 outline-none focus:border-blue-500" />
            </div>
            <div className="flex justify-end pt-4">
                <button onClick={onNext} disabled={!step1.senderName || !step1.receiverName} className="bg-blue-400 px-6 py-2 rounded font-bold hover:bg-blue-700 disabled:opacity-50">
                    Lanjut ke Step 2
                </button>
            </div>
        </div>
    );
}
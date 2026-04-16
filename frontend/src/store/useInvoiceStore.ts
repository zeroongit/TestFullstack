import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface InvoiceItem {
    id: number; 
    code: string;
    name: string; 
    quantity: number;
    price: number;
    subtotal: number;
}

interface InvoiceState {
    token: string | null;
    role: "Admin" | "Kerani" | null;
    step1: {
        senderName: string;
        senderAddress: string;
        receiverName: string;
        receiverAddress: string;
    };
    items: InvoiceItem[];
    setAuth: (token: string, role: "Admin" | "Kerani") => void;
    updateStep1: (data: InvoiceState['step1']) => void;
    addItem: (item: InvoiceItem) => void;
    setItems: (items: InvoiceItem[]) => void; // Tambahan untuk mempermudah update tabel
    resetForm: () => void;
    logout: () => void;
}

export const useInvoiceStore = create<InvoiceState>()(
    persist(
        (set) => ({
            token: null,
            role: null,
            step1: { senderName: "", senderAddress: "", receiverName: "", receiverAddress: "" },
            items: [],

            setAuth: (token, role) => set({ token, role }),
            
            updateStep1: (data) => set({ step1: data }),
            
            // Perbaikan addItem agar tidak duplikat (opsional tapi bagus untuk UX)
            addItem: (newItem) => set((state) => {
                const existingIndex = state.items.findIndex(i => i.code === newItem.code);
                if (existingIndex > -1) {
                    const updatedItems = [...state.items];
                    updatedItems[existingIndex].quantity += newItem.quantity;
                    updatedItems[existingIndex].subtotal = updatedItems[existingIndex].quantity * updatedItems[existingIndex].price;
                    return { items: updatedItems };
                }
                return { items: [...state.items, newItem] };
            }),

            setItems: (items) => set({ items }),

            resetForm: () => set({
                step1: { senderName: "", senderAddress: "", receiverName: "", receiverAddress: "" },
                items: []
            }),

            logout: () => set({
                token: null,
                role: null,
                step1: { senderName: "", senderAddress: "", receiverName: "", receiverAddress: "" },
                items: []
            }),
        }),
        {
            name: 'invoice-storage',
            // PENTING: Gunakan partialize jika ingin token tidak hilang saat reset form
            // tapi tetap tersimpan di storage terpisah.
        }
    )
);
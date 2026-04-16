import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface InvoiceState {
    token: string | null;
    role: "Admin" | "Kerani" | null;
    step1: {
        senderName: string;
        senderAddress: string;
        receiverName: string;
        receiverAddress: string;
    };
    items: Array<{
        id: number;
        code: string;
        name: string;
        quantity: number;
        price: number;
        subtotal: number;
    }>;
    setAuth: (token: string, role: "Admin" | "Kerani") => void;
    updateStep1: (data: InvoiceState['step1']) => void;
    addItem: (item: any) => void;
    resetForm: () => void;
}

export const useInvoiceStore = create<InvoiceState>()(
    persist(
        (set) => ({
            token: null,
            role: null,
            step1: {senderName: "", senderAddress: "", receiverName: "", receiverAddress: ""},
            items: [],

            setAuth: (token, role) => set({ token, role}),
            updateStep1: (data) => set({ step1: data }),
            addItem: (item) => set((state) => ({ items: [...state.items, item] })),
            resetForm: () => set({
                step1: {senderName: "", senderAddress: "", receiverName: "", receiverAddress: ""},
                items: []
            }),
        }),
        {
            name: 'invoice-storage',
        }
    )
);
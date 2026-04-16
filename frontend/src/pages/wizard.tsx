import { useState } from 'react';
import { useInvoiceStore } from '@/store/useInvoiceStore';
import Step1Form from '@/components/wizard/Step1Form';
import Step2Form from '@/components/wizard/Step2Form';
import Step3Review from '@/components/wizard/Step3Review';
import { usehHasHydrated } from '@/hooks/useHasHydrated';

export default function WizardPage() {
    const [currentStep, setCurrentStep] = useState(1);
    const hasHydrated = usehHasHydrated();


    if (!hasHydrated) return <div className='p-10 text-white'>Loading...</div>

    return (
        <div className='min-h-screen bg-slate-900 text-white p-4 md:p-10'>
            <div className='max-w-4xl mx-auto'>
                <div className='flex justify-between mb-8'>
                    {[1, 2, 3,].map((s) => (
                        <div key={s} className={`flex-1 text-center border-b-4 py-2 ${currentStep >= s ? 'border-blue-500 text-blue-400' : 'border-slate-700 text-slate-500'}`}>
                            Step {s}
                        </div>
                    ))}
                </div>

                <div className='bg-slate-800 p-6 rounded -xl border-slate-700 shadow-xl'>
                    {currentStep === 1 && <Step1Form onNext={() => setCurrentStep(2)} />}
                    {currentStep === 2 && <Step2Form onNext={() => setCurrentStep(3)} onPrev={() => setCurrentStep(1)} />}
                    {currentStep === 3 && <Step3Review onPrev={() => setCurrentStep(2)}/>}
                </div>
            </div>
        </div>
    );
}

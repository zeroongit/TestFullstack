import { useState, useEffect } from 'react';

export const usehHasHydrated = () => {
    const [hasHydrated, sethasHydrated] = useState(false);
    useEffect(() => {
        sethasHydrated(true);
    }, []);
    return hasHydrated;
};
import React, { useEffect, useState } from 'react';
import ArrivalsBoard from './ArrivalsBoard';
import './Platform.css'

function Platform({ line, stop_id, N = 10 }) {
    const [arrivals, setArrivals] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        let mounted = true;

        async function fetchArrivals(isPolling = false) {
            if (!isPolling) {
                setLoading(true);
                setError(null);
            }

            // Only delay on first load
            if (!isPolling) await new Promise(resolve => setTimeout(resolve, 1000));

            if (!mounted) return;

            try {
                const res = await fetch(`https://stand-clear.vercel.app/api/arrivals/?line=${line}&stop_id=${stop_id}&N=${N}`);
                if (!res.ok) {
                    throw new Error(`Server responded with status ${res.status}`);
                }
                const data = await res.json();
                if (mounted) {
                    setArrivals(Array.isArray(data) ? data : []);
                    if (!isPolling) setError(null);
                }
            } catch (err) {
                if (mounted) {
                    // On polling error, maybe keep old data? For now, standard error handling
                    setError(err.message);
                    setArrivals([]);
                }
            }
            if (mounted && !isPolling) setLoading(false);
        }

        fetchArrivals();
        const interval = setInterval(() => fetchArrivals(true), 15000);

        return () => {
            mounted = false;
            clearInterval(interval);
        };
    }, [line, stop_id, N]);

    if (loading) return <div className='loading-text'>Loading arrivals...</div>;
    if (error) return <div className='loading-text'>Error: {error}</div>;
    return <ArrivalsBoard arrivals={arrivals} />;
}

export default Platform;

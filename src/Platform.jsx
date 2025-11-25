import React, { useEffect, useState } from 'react';
import ArrivalsBoard from './ArrivalsBoard';
import './Platform.css'

function Platform({ line, stop_id, N = 10 }) {
    const [arrivals, setArrivals] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function fetchArrivals() {
            setLoading(true);
            setError(null);
            await new Promise(resolve => setTimeout(resolve, 1000));
            try {
                const res = await fetch(`https://stand-clear.vercel.app/api/arrivals/?line=${line}&stop_id=${stop_id}&N=${N}`);
                if (!res.ok) {
                    throw new Error(`Server responded with status ${res.status}`);
                }
                const data = await res.json();
                setArrivals(Array.isArray(data) ? data : []);
            } catch (err) {
                setError(err.message);
                setArrivals([]);
            }
            setLoading(false);
        }
        fetchArrivals();
    }, [line, stop_id, N]);

    if (loading) return <div className='loading-text'>Loading arrivals...</div>;
    if (error) return <div className='loading-text'>Error: {error}</div>;
    return <ArrivalsBoard arrivals={arrivals} />;
}

export default Platform;

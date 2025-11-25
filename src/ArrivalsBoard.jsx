import React from 'react';
import "./ArrivalsBoard.css";

function formatTime(seconds) {
    const s = Number(seconds);

    if (s <= 60) return "Arriving";
    return `${Math.floor(s / 60)} min`;
}

const ArrivalsBoard = ({ arrivals }) => {
    if (!arrivals || arrivals.length === 0) return null;
    const stopName = arrivals[0].friendly_stop;

    return (
        <div className="board">
            <div className="station-name">
                {stopName}
            </div>
            {arrivals.map((train, idx) => (
                <div
                    key={`${train.line}-${idx}`}
                    className={`row ${idx !== arrivals.length - 1 ? "with-divider" : ""}`}
                >
                    <div className="left">
                        <img
                            src={`/trains/${train.line}.svg`}
                            alt={`${train.line} train`}
                            className="line-badge"
                        />

                        <div className="direction-block">
                            <div className="direction">
                                {train.destination_direction}
                            </div>
                        </div>
                    </div>

                    <div
                        className={`time ${train.time_to_arrival <= 60 ? "arriving" : ""
                            }`}
                    >
                        {formatTime(train.time_to_arrival)}
                    </div>
                </div>
            ))}
        </div>)
};

export default ArrivalsBoard;

import { useState } from "react";

export default function Events() {
    const [eventsOpen, setEventsOpen] = useState(false);
    const [selectedEvents, setSelectedEvents] = useState([]);
    const events = ["media.pause", "media.stop", "media.play", "media.resume"];

    const toggleEvent = (event) => {
        setSelectedEvents(prev =>
            prev.includes(event)
                ? prev.filter(id => id !== event)
                : [...prev, event]
        );
    };

    const isEventDisabled = (event) => {
        const pauseStopEvents = ["media.pause", "media.stop"];
        const playResumeEvents = ["media.play", "media.resume"];

        if (pauseStopEvents.includes(event)) {
            // Disable pause/stop if play/resume is selected
            return selectedEvents.some(e => playResumeEvents.includes(e));
        }

        if (playResumeEvents.includes(event)) {
            // Disable play/resume if pause/stop is selected
            return selectedEvents.some(e => pauseStopEvents.includes(e));
        }

        return false;
    };

    return (
        <div
            className="select-events"
            onBlur={(e) => {
                if (!e.currentTarget.contains(e.relatedTarget)) {
                    setEventsOpen(false);
                }
            }}
            tabIndex={0}
        >
            <div
                className="select-events-header"
                onClick={() => setEventsOpen(prev => !prev)}
            >
                {selectedEvents.length === 0 ? "Select Events..." : `(${selectedEvents.length}) Event${selectedEvents.length === 1 ? "" : "s"} Selected`} {eventsOpen ? '◀' : '▶'}
            </div>
            {eventsOpen && (
                <div className="select-events-list">
                    {events.map(event => (
                        <label
                            key={event}
                            className="select-event-item"
                            style={{
                                opacity: isEventDisabled(event) ? 0.5 : 1,
                                cursor: isEventDisabled(event) ? 'not-allowed' : 'pointer'
                            }}
                        >
                            <input
                                type="checkbox"
                                checked={selectedEvents.includes(event)}
                                onChange={() => toggleEvent(event)}
                                disabled={isEventDisabled(event)}
                            />
                            {event.split('.')[1][0].toUpperCase() + event.split('.')[1].slice(1) + 's Media'}
                        </label>
                    ))}
                </div>
            )}
        </div>
    );
}

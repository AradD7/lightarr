import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import normalBulb from "/assets/bulbs/normalBulb.png"
import bulkyBulb from "/assets/bulbs/bulkyBulb.png"
import vintageBulb from "/assets/bulbs/vintageBulb.png"
import slimBulb from "/assets/bulbs/slimBulb.png"
import gu10Bulb from "/assets/bulbs/gu10Bulb.png"
import { useState } from "react";

const fetchAllBulbs = async () => {
    const response = await fetch("http://localhost:10100/api/bulbs");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
}

const fetchNewBulbs = async () => {
    const response = await fetch("http://localhost:10100/api/bulbs/refresh");
    if (!response.ok) throw new Error('Refresh bulbs response was not ok');
    return response.json();
}

const flashBulbRequest = async (mac) => {
    const response = await fetch("http://localhost:10100/api/bulbs/flash", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ mac }),
    });
    if (!response.ok) throw new Error('Flash bulb request failed');
};

const updateBulbName = async ({ mac, name }) => {
    const response = await fetch("http://localhost:10100/api/bulbs/name", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mac, name }),
    });
    if (!response.ok) throw new Error('Failed to update bulb name');
};

const updateBulbType = async ({ mac, type }) => {
    const response = await fetch("http://localhost:10100/api/bulbs/type", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mac, type }),
    });
    if (!response.ok) throw new Error('Failed to update bulb type');
};

const getBulbImg = (bulbType) => {
    const bulbTypeLower = bulbType.toLowerCase();
    if (bulbTypeLower === "normal") return normalBulb;
    if (bulbTypeLower === "bulky") return bulkyBulb;
    if (bulbTypeLower === "slim") return slimBulb;
    if (bulbTypeLower === "gu10") return gu10Bulb;
    if (bulbTypeLower === "vintage") return vintageBulb;
    return normalBulb;
}

export default function Bulbs() {
    const [flashingMac, setFlashingMac] = useState(new Set());
    const [changingBulbName, setChangingBulbName] = useState(null);
    const [newName, setNewName] = useState('');
    const [changingBulbType, setChangingBulbType] = useState(null);
    const [newType, setNewType] = useState('');


    const queryClient = useQueryClient();

    const { data: bulbs, isLoading, error } = useQuery({
        queryKey: ['bulbs'],
        queryFn: fetchAllBulbs,
    });

    const { data: numOfNewBulbs, refetch: refreshBulbs, isFetching } = useQuery({
        queryKey: ['numOfNewBulbs'],
        queryFn: fetchNewBulbs,
        enabled: false,
    });

    const flashBulbMutation = useMutation({
        mutationFn: flashBulbRequest,
    });

    const updateNameMutation = useMutation({
        mutationFn: updateBulbName,
        onSuccess: (data, variables) => {
            queryClient.setQueryData(['bulbs'], (oldBulbs) => {
                return oldBulbs.map(bulb =>
                    bulb.mac === variables.mac ?
                        { ...bulb, name: variables.name } :
                        bulb
                );
            });
        },
        onError: (error) => {
            console.log(error)
        },
    });

    const updateTypeMutation = useMutation({
        mutationFn: updateBulbType,
        onSuccess: (data, variables) => {
            queryClient.setQueryData(['bulbs'], (oldBulbs) => {
                return oldBulbs.map(bulb =>
                    bulb.mac === variables.mac ?
                        { ...bulb, type: variables.type } :
                        bulb
                );
            });
        },
    });

    const flashBulb = (mac) => {
        setFlashingMac(prev => new Set(prev).add(mac));
        flashBulbMutation.mutate(mac);

        setTimeout(() => {
            setFlashingMac(prev => {
                const newSet = new Set(prev);
                newSet.delete(mac);
                return newSet;
            });
        }, 2900);
    }

    const handleRefreshBulbs = async () => {
        const result = await refreshBulbs();
        if (result.data?.num > 0) {
            // Refetch all bulbs to include the new ones
            queryClient.invalidateQueries({ queryKey: ['bulbs'] });
        }
    };

    const changeBulbName = (mac, currentName) => {
        setChangingBulbName(mac);
        setNewName(currentName);
    };

    const handleNameSubmit = (mac) => {
        if (newName.trim()) {
            setChangingBulbName(null);
            updateNameMutation.mutate({ mac, name: newName.trim() });
        } else {
            setChangingBulbName(null);
        }
    };

    const changeBulbType = (mac, currentType) => {
        setChangingBulbType(mac);
        setNewType(currentType);
    };

    return (
        <div className="bulbs-page-contents">
            {isFetching ?
                <h1>
                    Looking for additional bulbs on the network...
                </h1> :
                numOfNewBulbs &&
                <h1>
                    {numOfNewBulbs.num === 0
                        ? "Found no new light bulbs on the network"
                        : `Found and added ${numOfNewBulbs.num} new ${numOfNewBulbs.num > 1 ? "bulbs" : "bulb"}!`}
                </h1>
            }
            <div className="bulbs-page">
                {error ? (
                    <h2 className="bulb-page-status">Something went wrong getting bulbs</h2>
                ) :
                    isLoading ?
                        <h2 className="bulb-page-status">Loading bulbs...</h2> :
                        bulbs?.length &&
                        bulbs?.slice()
                            .sort((a, b) => a.mac.localeCompare(b.mac))
                            .map(bulb =>
                                <section
                                    className="bulb-item" key={bulb.ip}
                                    style={{ opacity: bulb.isReachable ? 1 : 0.5 }}
                                >
                                    {changingBulbType === bulb.mac ? (
                                        <div
                                            className="bulb-type-picker"
                                            onBlur={() => setChangingBulbType(null)}
                                            tabIndex={0}
                                        >
                                            {['normal', 'bulky', 'slim', 'gu10', 'vintage'].map(type => (
                                                <img
                                                    key={type}
                                                    src={getBulbImg(type)}
                                                    alt={`${type} bulb`}
                                                    onClick={() => {
                                                        setChangingBulbType(null);
                                                        updateTypeMutation.mutate({ mac: bulb.mac, type });
                                                    }}
                                                    className={`bulb-type-option ${newType === type ? 'selected' : ''}`}
                                                />
                                            ))}
                                        </div>
                                    ) : (
                                        <img
                                            src={getBulbImg(bulb.type)}
                                            alt={`${bulb.type} wiz light bulb`}
                                            onClick={() => bulb.isReachable && changeBulbType(bulb.mac, bulb.type)}
                                            style={{ cursor: bulb.isReachable ? 'pointer' : 'default' }}
                                        />
                                    )}
                                    <section className="bulb-item-info">
                                        {changingBulbName === bulb.mac ? (
                                            <input
                                                type="text"
                                                value={newName}
                                                onChange={(e) => setNewName(e.target.value)}
                                                onKeyDown={(e) => {
                                                    if (e.key === 'Enter') handleNameSubmit(bulb.mac);
                                                    if (e.key === 'Escape') setChangingBulbName(null);
                                                }}
                                                onBlur={() => setChangingBulbName(null)}
                                                autoFocus
                                                className="bulb-name-input"
                                            />
                                        ) : (
                                            <h2 className="bulb-name">
                                                {bulb.name}
                                                {bulb.isReachable &&
                                                    <span
                                                        className="bulb-name-edit-icon material-symbols-outlined"
                                                        onClick={() => changeBulbName(bulb.mac, bulb.name)}
                                                    > edit
                                                    </span>
                                                }
                                            </h2>
                                        )}
                                        <h2 className="bulb-mac">
                                            {bulb.mac.toUpperCase()}
                                        </h2>
                                        <button
                                            className="bulb-flash-button"
                                            onClick={() => flashBulb(bulb.mac)}
                                            disabled={flashingMac.has(bulb.mac) || !bulb.isReachable}
                                        >
                                            {flashingMac === bulb.mac ? "Flashing..." : "Flash Light"}
                                        </button>
                                    </section>
                                </section>
                            )}
            </div>
            <button
                className="bulbs-refresh-button"
                onClick={() => handleRefreshBulbs()}
                disabled={isFetching}
            >
                {isFetching ? "Searching..." : "Add Bulbs"}
            </button>
        </div>
    )
}

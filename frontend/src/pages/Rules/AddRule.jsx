import { useQuery, useMutation } from "@tanstack/react-query";

import Accounts from "./Accounts";
import Devices from "./Devices";
import Events from "./Events";
import Commands from "./Commands";
import BulbsMac from "./BulbsMac";
import { useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

const fetchAccounts = async () => {
    const response = await fetch("http://localhost:10100/api/accounts");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const fetchDevices = async () => {
    const response = await fetch("http://localhost:10100/api/devices");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const fetchAllBulbs = async () => {
    const response = await fetch("http://localhost:10100/api/bulbs");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
}

const addRule = async ({ accounts, devices, events, commandBulbPairs }) => {
    const data = {
        condition: {
            event: events,
            account: accounts.map(acc => ({
                id: acc,
            })),
            device: devices.map(dev => ({
                id: dev,
            })),
        },
        action: commandBulbPairs.map(pair => {
            const action = {
                command: {
                    method: "setPilot",
                    params: {},
                },
                bulbsMac: pair.bulbsMacs,
            };

            pair.commands.forEach(cmd => {
                switch (cmd.command) {
                    case "Turn off":
                        action.command.params.state = false;
                        break;
                    case "Turn on":
                        action.command.params.state = true;
                        break;
                    case "Dim":
                        action.command.params.dimming = cmd.value;
                        break;
                    case "Change Color":
                        action.command.params.r = cmd.value.r;
                        action.command.params.g = cmd.value.g;
                        action.command.params.b = cmd.value.b;
                        break;
                    case "Change Temperature":
                        action.command.params.temp = cmd.value;
                        break;
                }
            });

            return action;
        })
    };

    const response = await fetch("http://localhost:10100/api/rules", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to add rule');
    return response.json();
}

export default function AddRule() {
    const navigate = useNavigate()
    const bottomRef = useRef(null);

    const [selectedAccounts, setSelectedAccounts] = useState([]);
    const [selectedDevices, setSelectedDevices] = useState([]);
    const [selectedEvents, setSelectedEvents] = useState([]);

    // Array of command/bulb pairs
    const [commandBulbPairs, setCommandBulbPairs] = useState([
        { id: 0, commands: [], bulbsMacs: [] }
    ]);
    const [nextPairId, setNextPairId] = useState(1);

    const [isEmpty, setIsEmpty] = useState(false);

    const { data: bulbs, isLoading: loadingBulbs, error: errorBulbs } = useQuery({
        queryKey: ['bulbs'],
        queryFn: fetchAllBulbs,
    });

    const { data: savedAccounts, isLoading: loadingSavedAccounts, error: errorSavedAccounts } = useQuery({
        queryKey: ['accounts'],
        queryFn: fetchAccounts,
    });

    const { data: savedDevices, isLoading: loadingSavedDevices, error: errorSavedDevices } = useQuery({
        queryKey: ['devices'],
        queryFn: fetchDevices,
    });

    const toggleAccount = (accountId) => {
        if (accountId === "toggleAll") {
            if (selectedAccounts.length === savedAccounts.length) {
                setSelectedAccounts([]);
            } else {
                setSelectedAccounts([...savedAccounts.map(acc => acc.id)]);
            }
        } else {
            setSelectedAccounts(prev =>
                prev.includes(accountId)
                    ? prev.filter(id => id !== accountId)
                    : [...prev, accountId]
            );
        }
    };

    const toggleDevice = (deviceId) => {
        if (deviceId === "toggleAll") {
            if (selectedDevices.length === savedDevices.length) {
                setSelectedDevices([]);
            } else {
                setSelectedDevices([...savedDevices.map(device => device.id)]);
            }
        } else {
            setSelectedDevices(prev =>
                prev.includes(deviceId)
                    ? prev.filter(id => id !== deviceId)
                    : [...prev, deviceId]
            );
        }
    };

    const toggleEvent = (event) => {
        setSelectedEvents(prev =>
            prev.includes(event)
                ? prev.filter(id => id !== event)
                : [...prev, event]
        );
    };

    const toggleCommand = (pairId, command) => {
        setCommandBulbPairs(prev => prev.map(pair => {
            if (pair.id !== pairId) return pair;

            const isSelected = pair.commands.some(c => c.command === command);

            if (isSelected) {
                return {
                    ...pair,
                    commands: pair.commands.filter(c => c.command !== command)
                };
            }

            // If selecting "Turn off", clear all other selections
            if (command === "Turn off") {
                return {
                    ...pair,
                    commands: [{ command: "Turn off" }]
                };
            }

            // If selecting anything else, remove "Turn off"
            let filtered = pair.commands.filter(c => c.command !== "Turn off");

            // Handle Change Color and Change Temperature mutual exclusivity
            if (command === "Change Color") {
                filtered = filtered.filter(c => c.command !== "Change Temperature");
            }
            if (command === "Change Temperature") {
                filtered = filtered.filter(c => c.command !== "Change Color");
            }

            // Add new command with default values
            let newCommand;
            if (command === "Dim") {
                newCommand = { command, value: 50 };
            } else if (command === "Change Color") {
                newCommand = { command, value: { r: 255, g: 255, b: 255 } };
            } else if (command === "Change Temperature") {
                newCommand = { command, value: 4000 };
            } else {
                newCommand = { command };
            }

            return {
                ...pair,
                commands: [...filtered, newCommand]
            };
        }));
    };

    const updateCommandValue = (pairId, command, value) => {
        setCommandBulbPairs(prev => prev.map(pair => {
            if (pair.id !== pairId) return pair;

            return {
                ...pair,
                commands: pair.commands.map(c =>
                    c.command === command ? { ...c, value } : c
                )
            };
        }));
    };

    const toggleBulbsMac = (pairId, bulbsMac) => {
        setCommandBulbPairs(prev => prev.map(pair => {
            if (pair.id !== pairId) return pair;

            if (bulbsMac === "toggleAll") {
                const availableBulbs = getAvailableBulbsForPair(pairId);
                if (pair.bulbsMacs.length === availableBulbs.length) {
                    return { ...pair, bulbsMacs: [] };
                } else {
                    return { ...pair, bulbsMacs: availableBulbs.map(b => b.mac) };
                }
            } else {
                return {
                    ...pair,
                    bulbsMacs: pair.bulbsMacs.includes(bulbsMac)
                        ? pair.bulbsMacs.filter(mac => mac !== bulbsMac)
                        : [...pair.bulbsMacs, bulbsMac]
                };
            }
        }));
    };

    const getAvailableBulbsForPair = (pairId) => {
        if (!bulbs) return [];

        const usedMacs = commandBulbPairs
            .filter(pair => pair.id !== pairId)
            .flatMap(pair => pair.bulbsMacs);

        return bulbs.filter(bulb => !usedMacs.includes(bulb.mac));
    };

    const handleAddCommandPair = () => {
        setCommandBulbPairs(prev => [
            ...prev,
            { id: nextPairId, commands: [], bulbsMacs: [] }
        ]);
        setNextPairId(prev => prev + 1);
        setTimeout(() => {
            bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
        }, 100);
    };

    const handleRemoveCommandPair = (pairId) => {
        if (commandBulbPairs.length === 1) return; // Keep at least one pair
        setCommandBulbPairs(prev => prev.filter(pair => pair.id !== pairId));
    };

    const addRuleMutation = useMutation({
        mutationFn: addRule,
        onSuccess: () => {
            navigate("/rules");
        }
    });

    const handleAddRule = () => {
        // Validate that all pairs have commands, events, and bulbs
        const hasEmptyPair = commandBulbPairs.some(pair =>
            pair.commands.length === 0 || pair.bulbsMacs.length === 0
        );

        if (hasEmptyPair || selectedEvents.length === 0) {
            setIsEmpty(true);
            return;
        }

        addRuleMutation.mutate({
            accounts: selectedAccounts,
            devices: selectedDevices,
            events: selectedEvents,
            commandBulbPairs: commandBulbPairs
        });
    }

    return (
        <>
            {isEmpty ? <h2>Cannot leave the following empty.</h2> : null}
            <div className="add-rule-page">
                <h1>When</h1>
                <section className="add-rule-page-line2">
                    <Accounts
                        data={savedAccounts}
                        addAccount={toggleAccount}
                        selectedAccounts={selectedAccounts}
                    />
                    <h1>,</h1>
                </section>
                <section onClick={() => setIsEmpty(false)}>
                    <Events
                        addEvent={toggleEvent}
                        isEmpty={isEmpty}
                    />
                </section>
                <h1>On</h1>
                <Devices
                    data={savedDevices}
                    addDevice={toggleDevice}
                    selectedDevices={selectedDevices}
                />
                <h1>Do</h1>
                {commandBulbPairs.map((pair, index) => (
                    <div key={pair.id} className="command-bulb-pair">
                        {index !== 0 && (
                            <h1>And Do</h1>
                        )}
                        <section onClick={() => setIsEmpty(false)}>
                            <Commands
                                addCommand={(command) => toggleCommand(pair.id, command)}
                                updateCommandValue={(command, value) => updateCommandValue(pair.id, command, value)}
                                isEmpty={isEmpty}
                                selectedCommands={pair.commands}
                            />
                        </section>
                        <h1>To</h1>
                        <section
                            onClick={() => setIsEmpty(false)}
                            className="add-rule-command-section"
                        >
                            <BulbsMac
                                data={getAvailableBulbsForPair(pair.id)}
                                addBulbsMac={(mac) => toggleBulbsMac(pair.id, mac)}
                                isEmpty={isEmpty}
                                selectedBulbsMacs={pair.bulbsMacs}
                            />
                            {commandBulbPairs.length > 1 && index !== 0 && (
                                <span
                                    className="material-symbols-outlined"
                                    onClick={() => handleRemoveCommandPair(pair.id)}
                                >
                                    delete
                                </span>
                            )}
                        </section>
                        {index === commandBulbPairs.length - 1
                            && getAvailableBulbsForPair(nextPairId + 1).length !== 0
                            && (
                                <span
                                    className="material-symbols-outlined"
                                    onClick={handleAddCommandPair}
                                >
                                    add_2
                                </span>
                            )}
                    </div>
                ))}
                <button
                    className="add-rule-button add-button"
                    onClick={handleAddRule}
                >
                    Add Rule
                </button>
                <div ref={bottomRef} />
            </div>
        </>
    )
}

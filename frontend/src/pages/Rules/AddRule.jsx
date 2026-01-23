import { useQuery } from "@tanstack/react-query";

import Accounts from "./Accounts";

import normalBulb from "/assets/bulbs/normalBulb.png"
import bulkyBulb from "/assets/bulbs/bulkyBulb.png"
import vintageBulb from "/assets/bulbs/vintageBulb.png"
import slimBulb from "/assets/bulbs/slimBulb.png"
import gu10Bulb from "/assets/bulbs/gu10Bulb.png"
import Devices from "./Devices";
import Events from "./Events";
import Commands from "./Commands";
import BulbsMac from "./BulbsMac";
import { useState } from "react";

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

export default function AddRule() {
    const [selectedAccounts, setSelectedAccounts] = useState([]);
    const [selectedDevices, setSelectedDevices] = useState([]);
    const [selectedEvents, setSelectedEvents] = useState([]);
    const [selectedCommands, setSelectedCommands] = useState([]);
    const [selectedBulbsMacs, setSelectedBulbsMacs] = useState([]);

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
        setSelectedAccounts(prev =>
            prev.includes(accountId)
                ? prev.filter(id => id !== accountId)
                : [...prev, accountId]
        );
    };
    const toggleDevice = (deviceId) => {
        setSelectedDevices(prev =>
            prev.includes(deviceId)
                ? prev.filter(id => id !== deviceId)
                : [...prev, deviceId]
        );
    };

    const toggleEvent = (event) => {
        setSelectedEvents(prev =>
            prev.includes(event)
                ? prev.filter(id => id !== event)
                : [...prev, event]
        );
    };

    const toggleCommand = (command) => {
        setSelectedCommands(prev => {
            const isSelected = prev.some(c => c.command === command);

            if (isSelected) {
                return prev.filter(c => c.command !== command);
            }

            // If selecting "Turn off", clear all other selections
            if (command === "Turn off") {
                return [{ command: "Turn off" }];
            }

            // If selecting anything else, remove "Turn off"
            let filtered = prev.filter(c => c.command !== "Turn off");

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

            return [...filtered, newCommand];
        });
    };

    const toggleBulbsMac = (bulbsMac) => {
        setSelectedBulbsMacs(prev =>
            prev.includes(bulbsMac)
                ? prev.filter(mac => mac !== bulbsMac)
                : [...prev, bulbsMac]
        );
    };

    const handleAddRule = () => {
        console.log(selectedAccounts, selectedEvents, selectedDevices, selectedCommands, selectedBulbsMacs)
    }

    return (
        <div className="add-rule-page">
            <h1>When</h1>
            <section className="add-rule-page-line2">
                <Accounts
                    data={savedAccounts}
                    addAccount={toggleAccount}
                />
                <h1>,</h1>
            </section>
            <Events
                addEvent={toggleEvent}
            />
            <h1>On</h1>
            <Devices
                data={savedDevices}
                addDevice={toggleDevice}
            />
            <h1>Do</h1>
            <Commands
                addCommand={toggleCommand}
            />
            <h1>To</h1>
            <BulbsMac
                data={bulbs}
                addBulbsMac={toggleBulbsMac}
            />
            <button
                className="add-rule-button add-button"
                onClick={handleAddRule}
            >
                Add Rule
            </button>
        </div>
    )
}

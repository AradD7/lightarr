import { useQuery } from "@tanstack/react-query";

import AddRule from "./AddRule";
import { useState } from "react";


const fetchRules = async () => {
    const response = await fetch("http://localhost:10100/api/rules");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

export default function Rules() {
    const [selectedAccounts, setSelectedAccounts] = useState([]);
    const [selectedDevices, setSelectedDevices] = useState([]);
    const [selectedEvents, setSelectedEvents] = useState([]);
    const [selectedCommands, setSelectedCommands] = useState([]);

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

    const { data: rules, isLoading: loadingRules, error: errorRules } = useQuery({
        queryKey: ['rules'],
        queryFn: fetchRules,
    });

    return (
        <div className="rules-page">
            <AddRule
                addAccount={toggleAccount}
                addDevice={toggleDevice}
                addEvent={toggleEvent}
                addCommand={toggleCommand}
            />
        </div>
    )
}

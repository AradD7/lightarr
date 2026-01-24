import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

import Accounts from "./Accounts";
import Devices from "./Devices";
import Events from "./Events";
import Commands from "./Commands";
import BulbsMac from "./BulbsMac";
import { useState } from "react";
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

const addRule = async ({ accounts, devices, events, commands, macs }) => {
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
    action: [{
      command: {
        method: "setPilot",
        params: {},
      },
      bulbsMac: macs,
    }]
  }
  commands.forEach(cmd => {
    switch (cmd.command) {
      case "Turn off":
        data.action[0].command.params.state = false;
        break;
      case "Turn on":
        data.action[0].command.params.state = true;
        break;
      case "Dim":
        data.action[0].command.params.dimming = cmd.value;
        break;
      case "Change Color":
        data.action[0].command.params.r = cmd.value.r;
        data.action[0].command.params.g = cmd.value.g;
        data.action[0].command.params.b = cmd.value.b;
        break;
      case "Change Temperature":
        data.action[0].command.params.temp = cmd.value;
        break;
    }
  });
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
  const [selectedAccounts, setSelectedAccounts] = useState([]);
  const [selectedDevices, setSelectedDevices] = useState([]);
  const [selectedEvents, setSelectedEvents] = useState([]);
  const [selectedCommands, setSelectedCommands] = useState([]);
  const [selectedBulbsMacs, setSelectedBulbsMacs] = useState([]);
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

  const addRuleMutation = useMutation({
    mutationFn: addRule,
    onSuccess: () => {
      console.log("Success!")
    }
  });

  const handleAddRule = () => {
    if (selectedCommands.length === 0 || selectedBulbsMacs.length === 0 || selectedEvents.length === 0) {
      setIsEmpty(true);
      return;
    }
    addRuleMutation.mutate({
      accounts: selectedAccounts,
      devices: selectedDevices,
      events: selectedEvents,
      commands: selectedCommands,
      macs: selectedBulbsMacs
    });
  }

  const handleAddCommand = () => {
    console.log("plus clicked")
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
        />
        <h1>Do</h1>
        <section onClick={() => setIsEmpty(false)}>
          <Commands
            addCommand={toggleCommand}
            isEmpty={isEmpty}
          />
        </section>
        <h1>To</h1>
        <section
          onClick={() => setIsEmpty(false)}
          className="add-rule-command-section"
        >
          <BulbsMac
            data={bulbs}
            addBulbsMac={toggleBulbsMac}
            isEmpty={isEmpty}
          />
          <span
            className="material-symbols-outlined"
            onClick={handleAddCommand}
          >
            add_2
          </span>
        </section>
        <button
          className="add-rule-button add-button"
          onClick={handleAddRule}
        >
          Add Rule
        </button>
      </div>
    </>
  )
}

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import tvLogo from "/assets/icons/tv.svg";
import webLogo from "/assets/icons/web.png";
import androidLogo from "/assets/icons/android.png";
import appleLogo from "/assets/icons/apple.png";
import defaultLogo from "/assets/icons/Plex.png";

const fetchAccounts = async () => {
    const response = await fetch("/api/accounts");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const fetchAllAccounts = async () => {
    const response = await fetch("/api/plex/accounts");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const fetchDevices = async () => {
    const response = await fetch("/api/devices");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const fetchAllDevices = async () => {
    const response = await fetch("/api/plex/devices");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const deleteAccount = async (id) => {
    const response = await fetch(`/api/accounts/${id}`, {
        method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete account');
    return response.json();
};

const addAccount = async (account) => {
    const response = await fetch("/api/accounts", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(account),
    });
    if (!response.ok) throw new Error('Failed to add account');
    return response.json();
};

const deleteDevice = async (id) => {
    const response = await fetch(`/api/devices/${id}`, {
        method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete device');
    return response.json();
};

const addDevice = async (device) => {
    const response = await fetch("/api/devices", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(device),
    });
    if (!response.ok) throw new Error('Failed to add device');
    return response.json();
};

// Helper function to determine device icon
const getDeviceIcon = (product) => {
    const productLower = product.toLowerCase();
    if (productLower.includes('tv')) return tvLogo;
    if (productLower.includes('android')) return androidLogo;
    if (productLower.includes('ios') || productLower.includes('iphone')) return appleLogo;
    if (productLower.includes('web')) return webLogo;
    return defaultLogo;
};

export default function PlexPage() {
    const [showAll, setShowAll] = useState(false);
    const [viewType, setViewType] = useState('accounts'); // 'accounts' or 'devices'
    const [dropdownOpen, setDropdownOpen] = useState(false);
    const [accountStatuses, setAccountStatuses] = useState({});

    const queryClient = useQueryClient();

    // Accounts queries
    const { data: savedAccounts, isLoading: loadingSavedAccounts, error: errorSavedAccounts } = useQuery({
        queryKey: ['accounts', 'saved'],
        queryFn: fetchAccounts,
        enabled: viewType === 'accounts',
    });

    const { data: allAccounts, isLoading: loadingAllAccounts, error: errorAllAccounts } = useQuery({
        queryKey: ['accounts', 'all'],
        queryFn: fetchAllAccounts,
        enabled: showAll && viewType === 'accounts',
    });

    // Devices queries
    const { data: savedDevices, isLoading: loadingSavedDevices, error: errorSavedDevices } = useQuery({
        queryKey: ['devices', 'saved'],
        queryFn: fetchDevices,
        enabled: viewType === 'devices',
    });

    const { data: allDevices, isLoading: loadingAllDevices, error: errorAllDevices } = useQuery({
        queryKey: ['devices', 'all'],
        queryFn: fetchAllDevices,
        enabled: showAll && viewType === 'devices',
    });

    // Account mutations
    const deleteAccountMutation = useMutation({
        mutationFn: deleteAccount,
        onMutate: (id) => {
            setAccountStatuses(prev => ({ ...prev, [id]: 'deleting' }));
        },
        onSuccess: (data, id) => {
            setAccountStatuses(prev => ({ ...prev, [id]: 'deleted' }));
            setTimeout(() => {
                queryClient.invalidateQueries({ queryKey: ['accounts', 'saved'] });
                queryClient.invalidateQueries({ queryKey: ['accounts', 'all'] });
                setAccountStatuses(prev => {
                    const newStatuses = { ...prev };
                    delete newStatuses[id];
                    return newStatuses;
                });
            }, 1000);
        },
        onError: (error, id) => {
            setAccountStatuses(prev => {
                const newStatuses = { ...prev };
                delete newStatuses[id];
                return newStatuses;
            });
            alert('Failed to delete account');
        }
    });

    const addAccountMutation = useMutation({
        mutationFn: addAccount,
        onMutate: (account) => {
            setAccountStatuses(prev => ({ ...prev, [account.id]: 'adding' }));
        },
        onSuccess: (data, account) => {
            setAccountStatuses(prev => ({ ...prev, [account.id]: 'added' }));
            setTimeout(() => {
                queryClient.invalidateQueries({ queryKey: ['accounts', 'saved'] });
                queryClient.invalidateQueries({ queryKey: ['accounts', 'all'] });
                setAccountStatuses(prev => {
                    const newStatuses = { ...prev };
                    delete newStatuses[account.id];
                    return newStatuses;
                });
            }, 1000);
        },
        onError: (error, account) => {
            setAccountStatuses(prev => {
                const newStatuses = { ...prev };
                delete newStatuses[account.id];
                return newStatuses;
            });
            alert('Failed to add account');
        }
    });

    // Device mutations
    const deleteDeviceMutation = useMutation({
        mutationFn: deleteDevice,
        onMutate: (id) => {
            setAccountStatuses(prev => ({ ...prev, [id]: 'deleting' }));
        },
        onSuccess: (data, id) => {
            setAccountStatuses(prev => ({ ...prev, [id]: 'deleted' }));
            setTimeout(() => {
                queryClient.invalidateQueries({ queryKey: ['devices', 'saved'] });
                queryClient.invalidateQueries({ queryKey: ['devices', 'all'] });
                setAccountStatuses(prev => {
                    const newStatuses = { ...prev };
                    delete newStatuses[id];
                    return newStatuses;
                });
            }, 1000);
        },
        onError: (error, id) => {
            setAccountStatuses(prev => {
                const newStatuses = { ...prev };
                delete newStatuses[id];
                return newStatuses;
            });
            alert('Failed to delete device');
        }
    });

    const addDeviceMutation = useMutation({
        mutationFn: addDevice,
        onMutate: (device) => {
            setAccountStatuses(prev => ({ ...prev, [device.id]: 'adding' }));
        },
        onSuccess: (data, device) => {
            setAccountStatuses(prev => ({ ...prev, [device.id]: 'added' }));
            setTimeout(() => {
                queryClient.invalidateQueries({ queryKey: ['devices', 'saved'] });
                queryClient.invalidateQueries({ queryKey: ['devices', 'all'] });
                setAccountStatuses(prev => {
                    const newStatuses = { ...prev };
                    delete newStatuses[device.id];
                    return newStatuses;
                });
            }, 1000);
        },
        onError: (error, device) => {
            setAccountStatuses(prev => {
                const newStatuses = { ...prev };
                delete newStatuses[device.id];
                return newStatuses;
            });
            alert('Failed to add device');
        }
    });

    const handleIconClick = (item) => {
        if (viewType === 'accounts') {
            if (showAll) {
                addAccountMutation.mutate(item);
            } else {
                deleteAccountMutation.mutate(item.id);
            }
        } else {
            if (showAll) {
                addDeviceMutation.mutate(item);
            } else {
                deleteDeviceMutation.mutate(item.id);
            }
        }
    };

    // Determine what to show based on viewType
    const savedIds = viewType === 'accounts'
        ? new Set(savedAccounts?.map(acc => acc.id) || [])
        : new Set(savedDevices?.map(dev => dev.id) || []);

    const unsavedItems = viewType === 'accounts'
        ? allAccounts?.filter(acc => !savedIds.has(acc.id)) || []
        : allDevices?.filter(dev => !savedIds.has(dev.id)) || [];

    const isLoading = viewType === 'accounts'
        ? (showAll ? loadingAllAccounts : loadingSavedAccounts)
        : (showAll ? loadingAllDevices : loadingSavedDevices);

    const error = viewType === 'accounts'
        ? (showAll ? errorAllAccounts : errorSavedAccounts)
        : (showAll ? errorAllDevices : errorSavedDevices);

    const itemsToShow = showAll ? unsavedItems : (viewType === 'accounts' ? savedAccounts : savedDevices);

    const getStatusText = (status) => {
        switch (status) {
            case 'deleting': return 'Deleting...';
            case 'deleted': return 'Deleted';
            case 'adding': return 'Adding...';
            case 'added': return 'Added';
            default: return '';
        }
    };

    const handleViewTypeChange = (type) => {
        setViewType(type);
        setDropdownOpen(false);
    };

    return (
        <div className="plex-page-accounts">
            <div className="header-with-dropdown">
                <h1
                    onClick={() => setDropdownOpen(!dropdownOpen)}
                >
                    {showAll
                        ? (viewType === 'accounts' ? "All Accounts" : "All Devices")
                        : (viewType === 'accounts' ? "Saved Accounts" : "Saved Devices")
                    }
                    <span
                        className="material-symbols-outlined dropdown-arrow"
                    >
                        {dropdownOpen ? 'arrow_drop_up' : 'arrow_drop_down'}
                    </span>
                </h1>

                {dropdownOpen && (
                    <div className="dropdown-menu">
                        <div
                            className={`dropdown-item ${viewType === 'accounts' ? 'active' : ''}`}
                            onClick={() => handleViewTypeChange('accounts')}
                        >
                            {`${showAll ? "All" : "Saved"} Accounts`}
                        </div>
                        <div
                            className={`dropdown-item ${viewType === 'devices' ? 'active' : ''}`}
                            onClick={() => handleViewTypeChange('devices')}
                        >
                            {`${showAll ? "All" : "Saved"} Devices`}
                        </div>
                    </div>
                )}
            </div>

            <div className="accounts">
                {error ? (
                    <h2>Something went wrong</h2>
                ) : isLoading ? (
                    <h2>Loading {viewType}...</h2>
                ) : itemsToShow?.length ? (
                    itemsToShow.map(item => (
                        <section className="account-item" key={item.id}>
                            {accountStatuses[item.id] && (
                                <div className="account-item-status">
                                    {getStatusText(accountStatuses[item.id])}
                                </div>
                            )}

                            <span
                                className={`material-symbols-outlined ${showAll ? "account-add-icon" : "account-delete-icon"}`}
                                onClick={() => handleIconClick(item)}
                            >
                                {showAll ? "person_add" : "delete"}
                            </span>

                            {viewType === 'accounts' ? (
                                <>
                                    <img src={item.thumb} alt={`${item.title}'s profile picture`} />
                                    <h1>{item.title}</h1>
                                </>
                            ) : (
                                <>
                                    <img src={getDeviceIcon(item.product)} alt={`${item.name} device icon`} />
                                    <h1>{item.name}</h1>
                                    <p className="device-product">{item.product}</p>
                                </>
                            )}
                        </section>
                    ))
                ) : (
                    <h2>
                        No {viewType} {showAll ? "available. All have been saved" : "saved"}.
                        {!showAll && ` Press below to add ${viewType}`}
                    </h2>
                )}
            </div>

            <button className="plex-add-button add-button" onClick={() => setShowAll(!showAll)}>
                {showAll
                    ? `View Saved ${viewType === 'accounts' ? 'Accounts' : 'Devices'}`
                    : `View All ${viewType === 'accounts' ? 'Accounts' : 'Devices'}`
                }
            </button>
        </div>
    );
}

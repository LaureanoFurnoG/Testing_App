import { createContext, useContext, useState } from "react";

type GroupsContextType = {
  refreshKey: number;
  refreshGroups: () => void;
};

const GroupsContext = createContext<GroupsContextType | null>(null);

export const GroupsProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [refreshKey, setRefreshKey] = useState(0);

  const refreshGroups = () => {
    setRefreshKey(prev => prev + 1);
  };

  return (
    <GroupsContext.Provider value={{ refreshKey, refreshGroups }}>
      {children}
    </GroupsContext.Provider>
  );
};

export const useGroups = () => {
  const ctx = useContext(GroupsContext);
  if (!ctx) throw new Error("useGroups must be used inside GroupsProvider");
  return ctx;
};

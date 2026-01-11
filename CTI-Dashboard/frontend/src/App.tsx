import { ThemeProvider } from "./components/theme-provider";
import { SidebarProvider, SidebarInset, SidebarTrigger } from "./components/ui/sidebar";
import { AppSidebar } from "./components/layout/sidebar";
import { Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import Threats from "./pages/Threats";
import Scrap from "./pages/Scrap";
import Forums from "./pages/Forums";
import { ModeToggle } from "./components/mode-toggle";


function App() {

  return (
      <div className="p-6 flex h-screen w-screen">
        <ThemeProvider>
          <ModeToggle />
          <SidebarProvider>
            <AppSidebar />
            <SidebarInset>
              <SidebarTrigger />
              <Routes>
                <Route path="/homepage" element={<HomePage />} />
                <Route path="/threats" element={<Threats />} />
                <Route path="/scrap" element={<Scrap />} />
                <Route path="/forums" element={<Forums />} />
              </Routes>
            </SidebarInset>
          </SidebarProvider>
        </ThemeProvider>
      </div>
  );
}

export default App;
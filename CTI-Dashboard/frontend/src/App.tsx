import { ThemeProvider } from "./components/theme-provider";
import { SidebarProvider, SidebarInset, SidebarTrigger } from "./components/ui/sidebar";
import { AppSidebar } from "./components/layout/sidebar";
import { Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import Threats from "./pages/Threats";
import Scrap from "./pages/Scrap";
import Forums from "./pages/Forums";
import { ModeToggle } from "./components/mode-toggle";
import { Toaster } from "./components/ui/sonner";


function App() {

  return (
    <ThemeProvider>
      <div className="p-6 flex h-screen w-screen">
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
          <Toaster />
      </div>
    </ThemeProvider>
  );
}

export default App;
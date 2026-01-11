import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Link } from "react-router-dom";

export function AppSidebar() {
  return (
    <Sidebar>
      <SidebarHeader className="text-xl font-bold">CTI Dashboard</SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarMenu>
            <SidebarMenuItem>
              <Link to="/homepage" className="w-full">
                Home Page
              </Link>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <Link to="/forums" className="w-full">
                Forums
              </Link>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <Link to="/threats" className="w-full">
                Threats
              </Link>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <Link to="/scrap" className="w-full">
                Scrap
              </Link>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter />
    </Sidebar>
  );
}

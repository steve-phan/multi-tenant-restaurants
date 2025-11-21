"use client";

import { useState } from 'react';
import api from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

interface Order {
  ID: string;
  Status: string;
  TotalAmount: number;
  CreatedAt: string;
}

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [restaurantId, setRestaurantId] = useState<string | null>(null);

  const fetchOrders = async () => {
    if (!restaurantId) return;
    try {
        const res = await api.get(`/restaurants/${restaurantId}/orders`);
        setOrders(res.data);
    } catch (e) {
        console.error(e);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Orders</h1>
      
      <div className="p-4 border rounded bg-yellow-50 mb-6">
        <Input 
            placeholder="Enter Restaurant ID" 
            value={restaurantId || ''} 
            onChange={(e) => setRestaurantId(e.target.value)} 
            className="max-w-sm mt-2"
        />
        <Button onClick={fetchOrders} className="mt-2">Load Orders</Button>
      </div>

      <Card>
        <CardHeader>
            <CardTitle>Recent Orders</CardTitle>
        </CardHeader>
        <CardContent>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Order ID</TableHead>
                        <TableHead>Time</TableHead>
                        <TableHead>Total</TableHead>
                        <TableHead>Status</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {orders.map(o => (
                        <TableRow key={o.ID}>
                            <TableCell>{o.ID.slice(0, 8)}...</TableCell>
                            <TableCell>{new Date(o.CreatedAt).toLocaleString()}</TableCell>
                            <TableCell>${o.TotalAmount}</TableCell>
                            <TableCell>{o.Status}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </CardContent>
      </Card>
    </div>
  );
}

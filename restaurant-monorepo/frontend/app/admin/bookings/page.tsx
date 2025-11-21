"use client";

import { useState } from 'react';
import api from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

interface Booking {
  ID: string;
  CustomerName: string;
  StartTime: string;
  NumberOfGuests: number;
  Status: string;
}

export default function BookingsPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [restaurantId, setRestaurantId] = useState<string | null>(null);

  const fetchBookings = async () => {
    if (!restaurantId) return;
    try {
        const res = await api.get(`/restaurants/${restaurantId}/bookings`);
        setBookings(res.data);
    } catch (e) {
        console.error(e);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Bookings</h1>
      
      <div className="p-4 border rounded bg-yellow-50 mb-6">
        <Input 
            placeholder="Enter Restaurant ID" 
            value={restaurantId || ''} 
            onChange={(e) => setRestaurantId(e.target.value)} 
            className="max-w-sm mt-2"
        />
        <Button onClick={fetchBookings} className="mt-2">Load Bookings</Button>
      </div>

      <Card>
        <CardHeader>
            <CardTitle>Upcoming Bookings</CardTitle>
        </CardHeader>
        <CardContent>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Customer</TableHead>
                        <TableHead>Time</TableHead>
                        <TableHead>Guests</TableHead>
                        <TableHead>Status</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {bookings.map(b => (
                        <TableRow key={b.ID}>
                            <TableCell>{b.CustomerName || 'Guest'}</TableCell>
                            <TableCell>{new Date(b.StartTime).toLocaleString()}</TableCell>
                            <TableCell>{b.NumberOfGuests}</TableCell>
                            <TableCell>{b.Status}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </CardContent>
      </Card>
    </div>
  );
}
